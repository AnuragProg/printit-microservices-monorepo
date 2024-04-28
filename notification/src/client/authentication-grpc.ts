import { auth_grpc } from '../proto_gen/authentication/auth';
import * as grpc from '@grpc/grpc-js';


const AUTH_GRPC_URI = process.env.AUTH_GRPC_URI || "localhost:50051";


class AuthGrpcClient {

	private authGrpcClient: auth_grpc.AuthenticationClient;

	constructor(){
		const credentials = grpc.credentials.createInsecure()
		this.authGrpcClient = new auth_grpc.AuthenticationClient(
			AUTH_GRPC_URI,
			credentials,
		);
	}


	async verifyToken(token: string): Promise<auth_grpc.User>{
		const payload = new auth_grpc.Token({token});
		return new Promise((res,rej)=>{
			this.authGrpcClient.VerifyToken(
				payload,
				(err, user)=>{
					if(err){
						rej(err);
					}
					res(user!);
				}
			);
		});
	}
}


export default new AuthGrpcClient();
export {
	AuthGrpcClient
}

