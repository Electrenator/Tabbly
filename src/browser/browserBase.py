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
        raise NotImplementedError()
