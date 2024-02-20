import random


def generate_otp(digits: int):
    low, high = 10**(digits-1), (10**digits)-1
    return random.randint(low, high)
