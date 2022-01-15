from browser.browserBase import BrowserBase
import lz4.block
import json


class Firefox(BrowserBase):
    def __init__(this):
        this.application_name = "GeckoMain"
        this.possible_tab_locations = [
            "~/.mozilla/firefox*/*.default/sessionstore-backups/recovery.jsonlz4"  # Firefox GNU/Linux
        ]

    def isRunning(this) -> bool:
        return super().isRunning()

    def getWindows(this) -> list:
        return super().getWindows()

    def getTabs(this) -> list:
        return super().getTabs()

    def parse_session_file(this, file_path: str) -> list:
        with open(file_path, "rb") as file:
            if file_path.find("firefox") != -1:
                file.read(8)  # ignore first firefox ID b"mozLz40\0"

            file_data = file.read()

            if file_path.endswith("lz4"):
                file_data = lz4.block.decompress(file_data).decode("utf-8")

            browser_data = json.loads(file_data)
            window_data = browser_data.get("windows")

            # browser_data json in at least following format from here -> {"windows": [{"tabs": [...]}]}
            return window_data if window_data != None else []
