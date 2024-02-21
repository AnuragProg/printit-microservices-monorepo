import asyncio
from typing import Literal
from fastapi import APIRouter, HTTPException
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from src.errors.not_found import NotFound
from src.model.user import UserModel, UserType, Purpose
from src.util.jwt_generator import generate_jwt
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

class ForgotPasswordDetails(BaseModel):
    email: str
    new_password: str

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
        async def forgot_password(fpd: ForgotPasswordDetails):
            return await self.forgot_password(fpd)

        @self.router.post('/verify-otp/{otp}')
        async def verify_otp(otp: int):
            return await self.verify_otp(otp)



    async def login(self, login_details: LoginDetails):
        try:
            user = await self.user_model.get_user_from_db(
                login_details.email,
                login_details.password
            )
            token = generate_jwt(
                _id=user._id
            )
            return JSONResponse(
                content={
                    'message': 'Logged in successfully',
                    'token': token
                },
                status_code=200
            )
        except NotFound:
            return JSONResponse(
                content={
                    'message': 'User not found'
                },
                status_code=404
            )
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

                # Check whether user exists
                await self.user_model.get_user_from_db(
                    email   =signup_details.email,
                    password=signup_details.password
                )
                return JSONResponse(
                    content={
                        'message': 'User already exists'
                    },
                    status_code=400
                )
            except NotFound:

                # generate otp
                otp = generate_otp(OTP_LEN)

                # cache user details
                await self.user_model.cache_for_signup(
                    key=str(otp),
                    ttl=OTP_TTL,
                    name=signup_details.name,
                    email=signup_details.email,
                    password=signup_details.password,
                    user_type=UserType.get(signup_details.user_type),
                )

                # send otp
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

    async def forgot_password(self, fpd: ForgotPasswordDetails):
        try:

            # check whether user exists
            await self.user_model.get_user_from_db(
                email=fpd.email,
                password=None
            )

            # generate otp
            otp = generate_otp(OTP_LEN)

            # cache user's details
            await self.user_model.cache_for_forgot_password(
                key=str(otp),
                ttl=OTP_TTL,
                email=fpd.email,
                new_password=fpd.new_password
            )

            # send otp
            await send_otp(otp, fpd.email)

            return JSONResponse(
                content={
                    'message': 'OTP sent successfully'
                },
                status_code=200
            )
        except NotFound:
            return JSONResponse(
                content={
                    'message': 'User not found'
                },
                status_code=404
            )
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

            # retrieve user details and delete it from cache
            info = await self.user_model.get_otp_info(str(otp))
            await self.user_model.del_otp_info(str(otp))

            if info['purpose'] == Purpose.SIGNUP.value:

                # save user details
                await self.user_model.save(
                    name=info['name'],
                    email=info['email'],
                    password=info['password'],
                    user_type=UserType.get(info['user_type'])
                )

                return JSONResponse(
                    content={
                        'message': 'User successfully signed up'
                    },
                    status_code=200
                )
            elif info['purpose'] == Purpose.FORGOT_PASSWORD.value:

                # update password
                await self.user_model.update_password(
                    email=info['email'],
                    new_password=info['new_password']
                )

                return JSONResponse(
                    content={
                        'message': 'Password updated successfully'
                    },
                    status_code=200
                )
            print(f"Found invalid purpose {info['purpose']}")
            raise Exception('Invalid purpose')
        except Exception as e:
            print(e)
            return JSONResponse(
                content={
                    'message': 'Something went wrong'
                },
                status_code=500
            )


