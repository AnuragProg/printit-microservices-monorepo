


class NotFound(Exception):
    def __init__(self, msg: str, resource: str):
        super().__init__(msg)
        self.resource = resource
