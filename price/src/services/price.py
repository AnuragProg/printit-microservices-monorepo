from uuid import UUID
from bson import ObjectId
from grpc import StatusCode
from grpc.aio import ServicerContext
from errors.not_found import NotFound
from proto_gen.price_pb2_grpc import PriceServicer
from proto_gen.price_pb2 import Empty, PriceIdAndShopId, PriceInfo, PageSize as ProtoPageSize, PageColor as ProtoPageColor
from model.price import get_price_model, PriceColor as PageColor, PricePageSize as PageSize


def str_to_proto_page_size(page_size: PageSize) -> ProtoPageSize:
    match page_size:
        case PageSize.A4:
            return ProtoPageSize.A4
    raise NotFound(msg="page size not found", resource="page size")


def str_to_proto_page_color(page_color: PageColor) -> ProtoPageColor:
    match page_color:
        case PageColor.BLACK_WHITE:
            return ProtoPageColor.BlackNWhite
        case PageColor.COLOR:
            return ProtoPageColor.Color

    raise NotFound(msg="page color not found", resource="page color")


class PriceService(PriceServicer):
    def __init__(self):
        pass

    async def HealthCheck(self, request: Empty, context: ServicerContext):
        return Empty()

    async def GetPriceInfoByPriceIdAndShopId(self, request: PriceIdAndShopId, context: ServicerContext):
        try:
            price_id = UUID(request.price_id)
            shop_id = ObjectId(request.shop_id)
            price_model = await get_price_model()
            price_info = await price_model.get_price_by_id_and_shop_id(price_id=price_id, shop_id=shop_id)

            return PriceInfo(
                _id=price_info._id.hex,
                shop_id=str(price_info.shop_id),
                shopkeeper_id=price_info.shopkeeper_id.hex,
                color=str_to_proto_page_color(price_info.color),
                page_size=str_to_proto_page_size(price_info.page_size),
                single_sided_price=price_info.single_sided_price,
                double_sided_price=price_info.double_sided_price
            )
        except NotFound as e:
            context.set_code(StatusCode.NOT_FOUND)
            context.set_details(f'{e.resource} not found')
        except Exception as e:
            print(e)
            context.set_code(StatusCode.INVALID_ARGUMENT)
            context.set_details("invalid shop id or price id")


