"""
Module to interact with browsers and gether every necessary browser detail to show
the usage count. Only the Browsers class it's getters need to be public for Tabbly
program to run.

TODO(Electrenator): Make a module class so every browser outputs data within
the same format. Not doing that will definitely be a problem when reading usage
from multiple browsers
"""
from abc import ABC, abstractmethod
import json
from psutil import process_iter, NoSuchProcess, AccessDenied, ZombieProcess
import lz4.block
import filesystem


class Browsers:
    """
    Class for getting the combined browser information from all open, and
    supported, browsers.
    """

    def count_tabs(this) -> int:
        """
        Getter for the detected tab count.
        """
        return len(_Firefox().get_tabs() if _Firefox().is_running() else [])

    def count_windows(this) -> int:
        """
        Getter for the detected window count.
        """
        return len(_Firefox().get_windows() if _Firefox().is_running() else [])


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
    def get_windows(this) -> list:
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
            for window in this.parse_session_file(file_path):
                browser_window_data.append(window)
        print(
            f"Read a total of {len(browser_window_data)} windows from '{this.__class__.__name__}'"
        )
        return browser_window_data

    @abstractmethod
    def get_tabs(this) -> list:
        """
        Returns:
            a list of all the active browser tabs within this browser. This is a
            concatenation of all the tabs open within all the windows.
        """
        browser_window_data = this.get_windows()
        browser_tab_data = []

        for window in browser_window_data:
            for tab in window.get("tabs"):
                browser_tab_data.append(tab)
        print(
            f"Read a total of {len(browser_tab_data)} tabs from '{this.__class__.__name__}'"
        )
        return browser_tab_data

    @abstractmethod
    def parse_session_file(this, file_path: str) -> list:
        """
        Parses a browsers session file within `this.possible_tab_locations` for windows
        and tabs. This is browser specific and should be written for every browser subclass.

        Returns:
            A object with the following structure::

                {
                    "windows": [
                        {
                            "tabs": [...]
                        }
                    ]
                }
        """
        raise NotImplementedError()


class _Firefox(_BrowserBase):
    """
    Firefox specific browser d
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
