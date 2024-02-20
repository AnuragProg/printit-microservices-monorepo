import uuid
import asyncio
from typing import Literal
from fastapi import APIRouter, HTTPException
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from src.errors.not_found import NotFound
from src.model.user import UserModel, UserType, Purpose
from src.util.otp_generator import generate_otp
from src.client.otp_service import send_otp


router = APIRouter()

'''
####needed routes
login
signup
forgot-password
verify
'''

OTP_LEN = 6
OTP_TTL = 10*60 # 10 mins

class LoginDetails(BaseModel):
    email: str
    password: str

class SignUpDetails(BaseModel):
    name: str
    email: str
    password: str
    user_type: Literal['customer', 'shopkeeper']


class UserRouter:
    def __init__(self, user_model: UserModel):
        self.router = APIRouter()
        self.user_model = user_model

    def setup_routes(self):

        @self.router.post('/login')
        async def login(login_details: LoginDetails):
            return await self.login(login_details)

        @self.router.post('/signup')
        async def signup(signup_details: SignUpDetails):
            return await self.signup(signup_details)

        @self.router.post('/forgot-password')
        async def forgot_password():
            return await self.forgot_password()

        @self.router.post('/verify-otp/{otp}')
        async def verify_otp(otp: int):
            return await self.verify_otp(otp)



    async def login(self, login_details: LoginDetails):
        try:
            user = await self.user_model.get_user_from_db(
                login_details.email,
                login_details.password
            )
            # user exists, create new jwt token
            return {'user': user}
        except NotFound:
            return {
                'message': 'User not found'
            }
        except Exception as e:
            print(e)
            return JSONResponse(
                content={
                    'message': 'Something went wrong'
                },
                status_code=500
            )

    async def signup(self, signup_details: SignUpDetails):
        try:
            try:
                await self.user_model.get_user_from_db(
                    email   =signup_details.email,
                    password=signup_details.password
                )
            except NotFound:
                otp = generate_otp(OTP_LEN)
                await self.user_model.cache(
                    key=str(otp),
                    ttl=OTP_TTL,
                    name=signup_details.name,
                    email=signup_details.email,
                    password=signup_details.password,
                    user_type=UserType.get(signup_details.user_type),
                    purpose=Purpose.SIGNUP
                )
                await send_otp(otp, signup_details.email)
                return {
                    'message': 'OTP sent successfully'
                }
        except Exception as e:
            print(e)
            return JSONResponse(
                content={
                    'message': 'Something went wrong'
                },
                status_code=500
            )

    async def forgot_password(self):
        try:
            pass
        except Exception as e:
            print(e)
            return JSONResponse(
                content={
                    'message': 'Something went wrong'
                },
                status_code=500
            )

    async def verify_otp(self, otp: int):
        try:
            info = await self.user_model.get_otp_info(str(otp))

            if info['purpose'] == Purpose.SIGNUP.value:
                await asyncio.gather(
                    self.user_model.del_otp_info(str(otp)),
                    self.user_model.save(
                        name=info['name'],
                        email=info['email'],
                        password=info['password'],
                        user_type=UserType.get(info['user_type'])
                    )
                )
                return JSONResponse(
                    content={
                        'message': 'User successfully signed up'
                    },
                    status_code=200
                )
            elif info['purpose'] == Purpose.FORGOT_PASSWORD.value:
                pass
            else:
                raise Exception('invalid purpose')

            return JSONResponse(
                content={
                    'message': 'everthing is ok'
                },
                status_code=200
            )
        except Exception as e:
            print(e)
            return JSONResponse(
                content={
                    'message': 'Something went wrong'
                },
                status_code=500
            )


