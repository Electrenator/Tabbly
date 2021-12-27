from pypresence import Presence


class DiscordPresence:
    def __init__(self, client_id: int):
        self._RPC = Presence(client_id)
        self._RPC.connect()
        print("Created Presence instance")

    def update(self, state: str):
        return self._RPC.update(state=state)

    def __del__(self):
        if (self._RPC):
            self._RPC.close()
        print("Closed Presence instance")
