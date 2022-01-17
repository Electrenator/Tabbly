from abc import ABC, abstractmethod
from psutil import process_iter, NoSuchProcess, AccessDenied, ZombieProcess
from filesystem import Filesystem


class BrowserBase(ABC):
    @abstractmethod
    def __init__(this):
        this.application_name = None
        this.possible_tab_locations = None

    @abstractmethod
    def isRunning(this) -> bool:
        """
        Checks if this browser is running by searching for the `application_name`
        within the active programs.

        Returns:
            true if the process is detected to be active, false if not.
        """
        if this.application_name == None:
            raise NotImplementedError()

        for process in process_iter():
            try:
                if this.application_name in process.name():
                    return True
            except (NoSuchProcess, AccessDenied, ZombieProcess):
                pass
        return False

    @abstractmethod
    def getWindows(this) -> list:
        """
        Returns:
            A list of active browser windows for this browser. A window usually
            has one or more `tabs` within itself.
        """
        if this.possible_tab_locations == None:
            raise NotImplementedError()

        session_files = Filesystem.find_files(this.possible_tab_locations[0])
        browser_window_data = []

        for file_path in session_files:
            for window in this.parse_session_file(file_path):
                browser_window_data.append(window)
        print(
            f"Read a total of {len(browser_window_data)} windows from '{this.__class__.__name__}'"
        )
        return browser_window_data

    @abstractmethod
    def getTabs(this) -> list:
        """
        Returns:
            a list of all the active browser tabs within this browser. This is a
            concatenation of all the tabs open within all the windows.
        """
        browser_window_data = this.getWindows()
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
