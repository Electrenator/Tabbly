"""
File for all the data used models within the Tabbly project.
"""


class BrowserData:
    """
    Saves the window tab data. Every list browser list entry is a window,
    while every value is the tab count.
    """

    def __init__(this, windows: list[int] | None = None):
        if not isinstance(windows, list) and windows is not None:
            raise TypeError("Type of input value should be an list or nothing")

        if isinstance(windows, list):
            for window in windows:
                if not isinstance(window, int):
                    raise TypeError("All list entries should be ints")

        this.windows = windows if windows is not None else []

    def add_window(this, window: int) -> None:
        """
        This function adds a window with it's total amount of tabs to the data array

        Args:
            window (int): Amount of tabs open within this window to add
        """
        if not isinstance(window, int):
            raise TypeError("Type of input value should be an int")

        this.windows.append(window)

    def get_data(this) -> list[int]:
        """
        Simple getter for all the browser data.

        Returns:
            list[int]: Returns a list of integers, where every entry is a window
            and each int value of this entry represents the total tabs open within
            this window.
        """
        return this.windows