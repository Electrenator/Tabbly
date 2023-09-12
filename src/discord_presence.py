"""
This module houses the Discord presence related communication class with it's
features.
"""
import sys
import time
import traceback
from pypresence import Presence
from pypresence.exceptions import InvalidID, DiscordError, DiscordNotFound

from models import Setting


class DiscordPresence:
    """
    Class handles the connection to Discord presence and it's maintenance when
    unable to connect or fully disconnected.
    """

    def __init__(this, client_id: int):
        this._max_retries = 3
        this.is_connected = False
        this._update_retries = 0
        this._presence_connection = None

        # Try to make a first connection
        try:
            this._presence_connection = Presence(client_id)
            if Setting.verbose:
                print("Created Presence instance")

            this.resume()
        except DiscordNotFound:
            print(
                "Discord doesn't seem to be open or installed, so I can't update "
                + "your status. Please install and open Discord so I can do my work "
                + "to the fullest potential. Limiting to file logging for now until "
                + "Discord is detected :)",
                file=sys.stderr,
            )
            this._update_retries = this._max_retries
            this._client_id = client_id

    def resume(this):
        """
        Resumes the connection to Discord Presence by reconnecting if not already
        connected. Does also handle known connectivity issues related to it and
        prints the unknown ones.
        Sets a connected class state when a connection was established.

        Will ignore resume requests when no presence connection was established
        to allow for file logging without presence working.
        """
        if this._presence_connection is None:
            try:
                this._presence_connection = Presence(this._client_id)
            except DiscordNotFound:
                print("Unable to detect Discord", file=sys.stderr)
                return

        if not this.is_connected:
            try:
                this._presence_connection.connect()
                this.is_connected = True
                this._update_retries = 0

                if Setting.verbose:
                    print("Connected to Presence")

            except (ConnectionRefusedError, ConnectionResetError):
                print(
                    "Unable to connect to Presence (Discord is probably closed)",
                    file=sys.stderr,
                )
            except DiscordError as ex:
                # Print and pass. Mostly being thrown on temporary discord states
                # like with not being logged in etc. Should be handles when new
                # ones are found though
                traceback.print_exception(
                    type(ex), ex, ex.__traceback__, file=sys.stderr
                )

    def update(this, state: str):
        """
        Updates the presence to have new contents. This may take some time to
        update within the Discord client since it usually does not update
        immediately when spamming the presence api with multiple updates per second.

        If there ever is a connection error while updating, it retries the class
        specified amount of times before giving up and setting the class to a
        disconnected state. Except for when no connection was established while
        starting Tabbly, in which case it will not retry.
        """
        seconds_between_retry = 1

        if this.is_connected and not Setting.dry_run:
            try:
                return this._presence_connection.update(state=state)
            except InvalidID:
                this.is_connected = False
                print("Presence suddenly disconnected", file=sys.stderr)

                # Retry adding status
                return this.update(state)

        if this._update_retries >= this._max_retries or Setting.dry_run:
            return None

        # Try connecting until this._max_retries
        time.sleep(seconds_between_retry)
        this._update_retries += 1
        print(
            (
                "Trying to reconnect "
                + f"({this._update_retries} of {this._max_retries} tries)..."
                if this._update_retries != 0
                else "..."
            ),
            file=sys.stderr,
        )
        this.resume()
        return this.update(state) if this._update_retries <= this._max_retries else None

    def pause(this):
        """
        Disconnects from the presence API and sets the class disconnected state.
        """
        if this.is_connected:
            this._presence_connection.close()
            this.is_connected = False

            if Setting.verbose:
                print("Closed connection to Presence")

    def __del__(this):
        if this.is_connected:
            this.pause()

        if Setting.verbose:
            print("Removed Presence instance")
