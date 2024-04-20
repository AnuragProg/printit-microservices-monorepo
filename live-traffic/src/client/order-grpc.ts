
import RequestInProgress from '../error/request-in-progress';
import { order_grpc } from '../proto_gen/order/order';
import * as grpc from '@grpc/grpc-js';


const ORDER_GRPC_URI = process.env.ORDER_GRPC_URI || "localhost:50055";


class OrderGrpcClient {

	private orderGrpcClient: order_grpc.OrderClient;
	private activeShopTrafficRequests: Set<string>;

	constructor(){
		const credentials = grpc.credentials.createInsecure()
		this.activeShopTrafficRequests = new Set();
		this.orderGrpcClient = new order_grpc.OrderClient(
			ORDER_GRPC_URI,
			credentials,
		);
	}

	getShopTraffic(
		shopId: string,
		updatedOnOrBeforeEpochMS: number,
	): Promise<{shopId: string, traffic: number}>{

		if(this.activeShopTrafficRequests.has(shopId)){
			throw new RequestInProgress();
		}

		const reqData = new order_grpc.GetShopTrafficRequest({
			shop_id: shopId,
			updated_on_or_before_epoch_ms: updatedOnOrBeforeEpochMS,
		});

		// add to active requests
		this.activeShopTrafficRequests.add(shopId);

		return new Promise((resolve, reject) => {
			this.orderGrpcClient.GetShopTraffic(
				reqData,
				(err, response)=>{
					// remove from active requests
					this.activeShopTrafficRequests.delete(shopId);
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
export {
	OrderGrpcClient
}
