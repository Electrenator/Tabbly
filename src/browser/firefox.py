import json
import lz4.block
from browser.browser_base import BrowserBase



class Firefox(BrowserBase):
    def __init__(this):
        super().__init__()
        this.possible_application_names = [
            # Firefox GNU/Linux (Ubuntu & Fedora tested)
            "GeckoMain",
            # Firefox GNU/Linux (Manjaro tested)
            "firefox"
        ]
        this.possible_tab_locations = [
            # Firefox GNU/Linux (Ubuntu, Fedora & Manjaro tested)
            "~/.mozilla/firefox*/*.default*/sessionstore-backups/recovery.jsonlz4",
        ]

    def is_running(this) -> bool:
        return super().is_running()

    def get_windows(this) -> list:
        return super().get_windows()

    def get_tabs(this) -> list:
        return super().get_tabs()

    def parse_session_file(this, file_path: str) -> list:
        with open(file_path, "rb") as file:
            if file_path.find("firefox") != -1:
                file.read(8)  # ignore first firefox ID b"mozLz40\0"

            file_data = file.read()

            if file_path.endswith("lz4"):
                file_data = lz4.block.decompress(file_data).decode("utf-8")

            browser_data = json.loads(file_data)
            window_data = browser_data.get("windows")

            # browser_data json in at least following format from here
            # -> {"windows": [{"tabs": [...]}]}
            return window_data if window_data is not None else []
