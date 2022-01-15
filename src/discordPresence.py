from pypresence import Presence


class DiscordPresence:
    def __init__(this, client_id: int):
        this._RPC = Presence(client_id)
        this._isConnected = False
        print("Created Presence instance")

        this.resume()

    def resume(this):
        if not this._isConnected:
            this._RPC.connect()
            this._isConnected = True
            print("Connected toPresence")

    def update(this, state: str):
        if this._isConnected:
            return this._RPC.update(state=state)

    def pause(this):
        if this._isConnected:
            this._RPC.close()
            this._isConnected = False
            print("Closed connection to Presence")

    def __del__(this):
        if this._isConnected:
            this.pause()
        print("Closed Presence instance")
