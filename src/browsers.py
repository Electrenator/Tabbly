"""
Module to interact with browsers and gether every necessary browser detail to show
the usage count. Only the Browsers class it's getters need to be public for Tabbly
program to run.
"""
from abc import ABC, abstractmethod
import json
from psutil import process_iter, NoSuchProcess, AccessDenied, ZombieProcess
import lz4.block
import filesystem
from models import BrowserData


class Browsers:
    """
    Class for getting the combined browser information from all open, and
    supported, browsers.
    """

    def get_windows(this) -> list[int]:
        """
        Getter for getting the total tabs per window

        Returns:
            list[int]: Returns a list where every entry is a window and the
            associated value is the tab count
        """
        return _Firefox().get_windows() if _Firefox().is_running() else []


class _BrowserBase(ABC):
    """
    Default class with functions and default that all browser specific classes
    should have.
    """

    def __init__(this):
        this.possible_application_names = None
        this.possible_tab_locations = None

    @abstractmethod
    def is_running(this) -> bool:
        """
        Checks if this browser is running by searching for the `possible_application_names`
        within the active programs.

        Returns:
            true if the process is detected to be active, false if not.
        """
        if this.possible_application_names is None:
            raise NotImplementedError()

        for process in process_iter():
            try:
                for name in this.possible_application_names:
                    if name in process.name():
                        return True
            except (NoSuchProcess, AccessDenied, ZombieProcess):
                pass
        return False

    @abstractmethod
    def get_windows(this) -> list[int]:
        """
        Returns:
            A list of active browser windows for this browser. A window usually
            has one or more `tabs` within itself.
        """
        if this.possible_tab_locations is None:
            raise NotImplementedError()

        session_files = filesystem.find_files(this.possible_tab_locations[0])
        browser_window_data = []

        for file_path in session_files:
            browser_window_data = this.parse_session_file(file_path).get_data()

        print(f"Read {browser_window_data} from '{this.__class__.__name__.replace('_', '')}'.")
        return browser_window_data

    @abstractmethod
    def parse_session_file(this, file_path: str) -> BrowserData:
        """
        Parses a browsers session file within `this.possible_tab_locations` for windows
        and tabs. This is browser specific and should be written for every browser subclass.

        Returns: An filled BrowserData object
        """
        raise NotImplementedError()


class _Firefox(_BrowserBase):
    """
    Firefox specific browser data
    """

    def __init__(this):
        super().__init__()
        this.possible_application_names = [
            # Firefox GNU/Linux (Ubuntu & Fedora tested)
            "GeckoMain",
            # Firefox GNU/Linux (Manjaro tested)
            "firefox",
        ]
        this.possible_tab_locations = [
            # Firefox GNU/Linux (Ubuntu, Fedora & Manjaro tested)
            "~/.mozilla/firefox*/*.default*/sessionstore-backups/recovery.jsonlz4",
        ]

    def is_running(this) -> bool:
        return super().is_running()

    def get_windows(this) -> list:
        return super().get_windows()

    def parse_session_file(this, file_path: str) -> BrowserData:
        raw_browser_data = ""
        browser_data = BrowserData()

        # Read and decode file data
        with open(file_path, "rb") as file:
            if file_path.find("firefox") != -1:
                file.read(8)  # ignore first firefox ID b"mozLz40\0"

            # Read and decompress file
            file_data_raw = file.read()
            file_data = lz4.block.decompress(file_data_raw).decode("utf-8")

            # Load inner json
            raw_browser_data = json.loads(file_data)

        # Read and insert window data into BrowserData object
        window_data = raw_browser_data.get("windows")
        for window in window_data:
            # Calculates tabs from the saved window object itself since the given
            # "Selected" value appears to be inaccurate with lots of tabs open.
            browser_data.add_window(len(window["tabs"]))

        return browser_data
