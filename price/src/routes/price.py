from uuid import UUID
from bson import ObjectId as oid
from grpc import RpcError
from pydantic import BaseModel
from typing import Optional
from fastapi import APIRouter, Depends
from fastapi.responses import JSONResponse
from errors.not_found import NotFound
from model.price import Price, PriceColor, PricePageSize
from proto_gen.auth_pb2 import User
from proto_gen.shop_pb2 import ShopInfo, ShopAndShopkeeperId
from middleware.auth import auth_shopkeeper, auth_user
from client.shop_grpc_client import get_shop_grpc_client
from model.price import get_price_model


router = APIRouter()

class CreatePriceModel(BaseModel):
    single_sided_price: float
    double_sided_price: float
    color: PriceColor
    page_size: PricePageSize

class UpdatePriceModel(BaseModel):
    single_sided_price: Optional[float] = None
    double_sided_price: Optional[float] = None
    color: Optional[PriceColor] = None
    page_size: Optional[PricePageSize] = None



# TODO: test this endpoint
@router.post('/shop/{shop_id}/prices')
async def create_price(
    shop_id: str,
    price_body: CreatePriceModel,
    user: User = Depends(auth_shopkeeper),
):
    try:
        # retrieve grpc client
        shop_grpc_client = await get_shop_grpc_client()
        shop : ShopInfo = await shop_grpc_client.GetShopByShopAndShopkeeperId(
            ShopAndShopkeeperId(
               shop_id=shop_id,
               shopkeeper_id=user._id
            )
        )

        # construct price
        price = Price(
            shop_id=oid(shop._id),
            shopkeeper_id=UUID(user._id),
            color=price_body.color,
            single_sided_price=price_body.single_sided_price,
            double_sided_price=price_body.double_sided_price,
            page_size=price_body.page_size
        )

        # retrieve price model client
        price_model = await get_price_model()

        # insert the price
        await price_model.insert_price(price)

        return JSONResponse(
            content={'message': 'price successfully created', 'price': price.to_dict()},
            status_code=200
        )
    except RpcError as e:
        print(e)
        return JSONResponse(
            content = {'message': 'shop not found'},
            status_code=404
        )
    except Exception as e:
        print(e)
        return JSONResponse(
            content = {'message': 'invalid request'},
            status_code=400
        )

# TODO: test this endpoint
@router.get('/shop/{shop_id}/prices')
async def get_prices(
    shop_id: str,
    user: User = Depends(auth_user)
):
    try:
        # check validity of shop_id
        if not oid.is_valid(shop_id):
            return JSONResponse(
                content = {'message': 'invalid shop id'},
                status_code=400
            )
        shop_id_oid = oid(shop_id)

        # retrieve price model client
        price_model = await get_price_model()

        # retrieve prices of the shop and convert them into dict
        prices = list(map(lambda price: price.to_dict(), await price_model.get_prices_by_shop_id(shop_id=shop_id_oid)))
        return JSONResponse(
            content = {'prices': prices},
            status_code=200
        )
    except NotFound as e:
        print(e)
        return JSONResponse(
            content = {'message': f'{e.resource} not found'},
            status_code=404
        )
    except Exception as e:
        print(e)
        return JSONResponse(
            content = {'message': 'invalid request'},
            status_code=400
        )

# TODO: test this endpoint
@router.get('/shop/{shop_id}/prices/{price_id}')
async def get_price(
    shop_id: str,
    price_id: str,
    user: User = Depends(auth_user)
):
    try:
        # check validity of shop_id
        if not oid.is_valid(shop_id):
            return JSONResponse(
                content = {'message': 'invalid shop id'},
                status_code=400
            )
        shop_id_oid = oid(shop_id)

        # retrieve price model client
        price_model = await get_price_model()

        # retrieve price
        price = await price_model.get_price_by_id(shop_id=shop_id_oid, price_id=UUID(price_id))

        return JSONResponse(
            content = {'price': price.to_dict()},
            status_code=200
        )
    except NotFound as e:
        return JSONResponse(
            content = {'message': f'{e.resource} not found'},
            status_code=404
        )
    except Exception:
        return JSONResponse(
            content = {'message': 'invalid request'},
            status_code=400
        )


# TODO: test this endpoint
@router.patch('/shop/{shop_id}/prices/{price_id}')
async def update_price(
    shop_id: str,
    price_id: str,
    price_body: UpdatePriceModel,
    user: User = Depends(auth_shopkeeper)
):
    try:
        shop_grpc_client = await get_shop_grpc_client()
        shop : ShopInfo = await shop_grpc_client.GetShopByShopAndShopkeeperId(
            ShopAndShopkeeperId(
               shop_id=shop_id,
               shopkeeper_id=user._id
            )
        )
        price_model = await get_price_model()
        await price_model.update_price(
            shop_id=oid(shop._id),
            price_id=UUID(price_id),
            single_sided_price=price_body.single_sided_price,
            double_sided_price=price_body.double_sided_price,
            color=price_body.color,
            page_size=price_body.page_size
        )
        return JSONResponse(
            content={'message': 'price updated successfully'},
            status_code=200
        )
    except RpcError as e:
        print(e)
        return JSONResponse(
            content={'message': 'shop not found'},
            status_code=404
        )
    except Exception as e:
        print(e)
        return JSONResponse(
            content={'message': 'invalid request'},
            status_code=400
        )

@router.delete('/shop/{shop_id}/prices/{price_id}')
async def delete_price(
    shop_id: str,
    price_id: str,
    user: User = Depends(auth_shopkeeper)
):
    try:
        shop_grpc_client = await get_shop_grpc_client()
        shop : ShopInfo = await shop_grpc_client.GetShopByShopAndShopkeeperId(
            ShopAndShopkeeperId(
               shop_id=shop_id,
               shopkeeper_id=user._id
            )
        )
        price_model = await get_price_model()
        await price_model.delete_price(price_id=UUID(price_id), shop_id=oid(shop_id))

        return JSONResponse(
            content={'message': 'price deleted successfully'},
            status_code=200
        )
    except RpcError as e:
        print(e)
        return JSONResponse(
            content={'message': 'shop not found'},
            status_code=404
        )
    except Exception as e:
        print(e)
        return JSONResponse(
            content={'message': 'invalid request'},
            status_code=400
        )
