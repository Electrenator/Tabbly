import sys
import time
import traceback
from pypresence import Presence
from pypresence.exceptions import InvalidID, DiscordError, DiscordNotFound


class DiscordPresence:
    def __init__(this, client_id: int):
        this._max_retries = 3
        this._is_connected = False
        this._update_retries = 0

        this._presence_connection = Presence(client_id)
        print("Created Presence instance")

        # Try to make a first connection
        try:
            this.resume()
        except DiscordNotFound as ex:
            traceback.print_exception(
                type(ex), ex, ex.__traceback__, file=sys.stderr)
            print(
                "Discord doesn't seem to be installed, so I can't update" +
                " your status. Please install Discord so I can do my work :)",
                file=sys.stderr
            )
            exit(-1)

    def resume(this):
        if not this._is_connected:
            try:
                this._presence_connection.connect()
                this._is_connected = True
                print("Connected to Presence")
                this._update_retries = 0
            except (ConnectionRefusedError, ConnectionResetError):
                print(
                    "Unable to connect to Presence (Discord is probably closed)",
                    file=sys.stderr
                )
            except DiscordError as ex:
                # Print and pass. Mostly being thrown on temporary discord states
                # like with not being logged in etc
                traceback.print_exception(
                    type(ex), ex, ex.__traceback__, file=sys.stderr
                )

    def update(this, state: str):
        seconds_between_retry = 1

        if this._is_connected:
            try:
                return this._presence_connection.update(state=state)
            except InvalidID:
                this._is_connected = False
                print("Presence suddenly disconnected", file=sys.stderr)

                # Retry adding status
                return this.update(state)

        if this._update_retries >= this._max_retries:
            return None

        # Try connecing until this._max_retries
        time.sleep(seconds_between_retry)
        this._update_retries += 1
        print(
            "Trying to reconnect " +
            f"({this._update_retries} of {this._max_retries} tries)" if this._update_retries != 0
            else "",
            file=sys.stderr
        )
        this.resume()
        return this.update(state) if this._update_retries <= this._max_retries else None

    def pause(this):
        if this._is_connected:
            this._presence_connection.close()
            this._is_connected = False
            print("Closed connection to Presence")

    def __del__(this):
        if this._is_connected:
            this.pause()
        print("Closed Presence instance")
