from pypresence import Presence


class DiscordPresence:
    def __init__(this, client_id: int):
        this._RPC = Presence(client_id)
        this._RPC.connect()
        print("Created Presence instance")

    def update(this, state: str):
        return this._RPC.update(state=state)
        
    def __del__(this):
        if this._RPC:
            this._RPC.close()
        print("Closed Presence instance")
