import glob
import os
import lz4.block
import json


class Browser:
    def countTabs(this) -> int:
        return len(this._get_tabs())

    def _get_tabs(this) -> list:
        possible_tab_locations = (
            "~/.mozilla/firefox*/*.default/sessionstore-backups/recovery.jsonlz4"
        )

        found_tab_files = this._find_files(possible_tab_locations)
        browser_tabs = []

        for file_path in found_tab_files:
            browser_windows = this._parse_browser_file(file_path)

            for window in browser_windows.get("windows"):
                tabs = window.get("tabs")
                print(f"Read a total of {len(tabs)} tabs from '{file_path}'")
                browser_tabs += tabs

        return browser_tabs

    def _find_files(this, path) -> list:
        """
        Finds files within given path while allowing relative linux like glob paths to be used.

        Replaces `~` with the users home directory, rest of the glob syntax is explained within the
        documentation; https://docs.python.org/3.8/library/glob.html
        """
        if path.startswith("~"):
            path = path.replace("~", os.path.expanduser("~"), 1)

        return glob.glob(path)

    def _parse_browser_file(this, file_path):
        with open(file_path, "rb") as file:
            if file_path.find("firefox") != -1:
                file.read(8)  # ignore first firefox ID b"mozLz40\0"

            file_data = file.read()

            if file_path.endswith("lz4"):
                file_data = lz4.block.decompress(file_data).decode("utf-8")

            browser_data = json.loads(file_data)

            # browser_data json in at least following format from here -> {"windows": [{"tabs": [...]}]}
            return browser_data
