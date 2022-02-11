from browser.firefox import Firefox


class Browsers:
    def count_tabs(this) -> int:
        return len(Firefox().get_tabs() if Firefox().is_running() else [])

    def count_windows(this) -> int:
        return len(Firefox().get_windows() if Firefox().is_running() else [])
