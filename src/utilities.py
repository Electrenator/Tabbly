from argparse import ArgumentParser


def init_parser() -> ArgumentParser:
    parser = ArgumentParser(
        description="A program for showing your browser tab usage within Discord's rich presence "
        + " and to log that usage.",
        epilog="Project source can be found over at; https://github.com/Electrenator/Tabbly",
        allow_abbrev=False,
    )
    parser.add_argument(
        "-v", "--verbose", action="store_true", help="show extended console output"
    )
    return parser
