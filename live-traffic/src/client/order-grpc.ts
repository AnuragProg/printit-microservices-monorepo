
import { order_grpc } from '../proto_gen/order/order';
import * as grpc from '@grpc/grpc-js';


const ORDER_GRPC_URI = process.env.ORDER_GRPC_URI || "localhost:50055";


class OrderGrpcClient {
	private orderGrpcClient: order_grpc.OrderClient;
	constructor(){
		const credentials = grpc.credentials.createInsecure()
		this.orderGrpcClient = new order_grpc.OrderClient(
			ORDER_GRPC_URI,
			credentials,
		);
	}


	getShopTraffic(
		shopId: string,
		updatedOnOrBefore: Date,
	): Promise<{shopId: string, traffic: number}>{
		const reqData = new order_grpc.GetShopTrafficRequest({
			shop_id: shopId,
			updated_on_or_before: updatedOnOrBefore.toISOString(),
		});
		return new Promise((resolve, reject) => {
			this.orderGrpcClient.GetShopTraffic(
				reqData,
				(err, response)=>{
					if(err){
						reject(err.message);
					}
					resolve({
						shopId: response!.shop_id,
						traffic: response!.traffic,
					});
				}
			);
		});
	}
}


export default new OrderGrpcClient();
