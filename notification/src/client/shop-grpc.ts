import { shop_grpc } from '../proto_gen/shop/shop';
import * as grpc from '@grpc/grpc-js';


const SHOP_GRPC_URI = process.env.SHOP_GRPC_URI || "localhost:50053";


class ShopGrpcClient {

	private shopGrpcClient: shop_grpc.ShopClient;

	constructor(){
		const credentials = grpc.credentials.createInsecure()
		this.shopGrpcClient = new shop_grpc.ShopClient(
			SHOP_GRPC_URI,
			credentials,
		);
	}

	async getShopById(shopId: string): Promise<shop_grpc.ShopInfo>{
		const payload = new shop_grpc.ShopId({
			_id: shopId,
		});
		return new Promise((res,rej)=>{
			this.shopGrpcClient.GetShopById(
				payload,
				(err, shopInfo)=>{
					if(err){
						rej(err);
					}
					res(shopInfo!);
				}
			);
		});
	}
}


export default new ShopGrpcClient();
export {
	ShopGrpcClient
}

