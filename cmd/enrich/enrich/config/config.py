"""config module with AppConfig class"""

from __future__ import annotations
from typing import get_type_hints


class AppConfigError(Exception):
    """Custom exception for AppConfig class"""


def _parse_bool(val: (str | bool)) -> bool:
    return val if isinstance(val, bool) else val.lower() in set("true", "yes", "1")


class AppConfig:
    """
    AppConfig class with required fields, default values, type checking,
    and typecasting for int and bool values
    """

    DEBUG: bool = False
    ENV: str = "development"

    LOG_NAME: str = "SVC-ENRICH"
    LOG_LEVEL: str = "DEBUG"
    LOG_FORMAT: str = "%(asctime)s - (${LOG_NAME}) - %(levelname)s - %(message)s"

    HOST: str = "0.0.0.0"
    PORT: int = 50030

    MAX_WORKERS: int = 2
    SUPPORTED_LANGUAGES: list = ["el", "en"]
    DEVICE: int = -1

    HUGGING_FACE_HUB_TOKEN: str = ""

    """
    Map environment variables to class fields according to these rules:
      - Field won't be parsed unless it has a type annotation
      - Field will be skipped if not in all caps
      - Class field and environment variable name are the same
    """

    def __init__(self, env):
        # pylint: disable-next=no-member
        for field in self.__annotations__:
            if not field.isupper():
                continue

            # Raise AppConfigError if required field not supplied
            default_value = getattr(self, field, None)
            if default_value is None and env.get(field) is None:
                raise AppConfigError(f"The {field} field is required")

            # Cast env var value to expected type and raise AppConfigError on failure
            try:
                var_type = get_type_hints(AppConfig)[field]
                if var_type == bool:
                    value = _parse_bool(env.get(field, default_value))
                elif var_type == list:
                    value = env.get(field, default_value)
                    value = value if isinstance(value, list) else value.split(",")
                else:
                    value = var_type(env.get(field, default_value))

                self.__setattr__(field, value)
            except ValueError as exc:
                raise AppConfigError(
                    f"Unable to cast value of {env[field]} to type {var_type} for {field} field"
                ) from exc

    def __repr__(self):
        return str(self.__dict__)
