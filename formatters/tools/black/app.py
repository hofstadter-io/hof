import os
import json
import base64
import lzma
import black
import urllib
import tempfile

from black.mode import TargetVersion
from flask import Flask, render_template, request, jsonify
from flask_cors import cross_origin

TEMP_DIR = tempfile.gettempdir()

BASE_URL = "https://black.vercel.app"
BLACK_VERSION = os.getenv("BLACK_VERSION")

TARGET_VERSIONS = {v.name.lower(): v for v in TargetVersion}


def get_black_version():
    lockfile = json.load(open("./Pipfile.lock"))
    package = lockfile["default"]["black"]
    version = package.get("version")
    ref = package.get("ref")

    if version:
        return version.lstrip("==")

    if ref:
        return ref[0:6]


app = Flask(__name__)
black_version = get_black_version()


def compress_state(data):
    compressed = lzma.compress(json.dumps(data).encode("utf-8"))
    return base64.urlsafe_b64encode(compressed).decode("utf-8")


def decompress_state(state):
    compressed = base64.urlsafe_b64decode(state)
    return json.loads(lzma.decompress(compressed))


def normalize_exception(exc):
    exception_str = f"{exc}"

    # Try to load contents dumped to tmp file.
    if "helpful: " in exception_str:
        try:
            _, file_path = exception_str.split("helpful: ")

            if file_path.startswith(TEMP_DIR):
                with open(file_path) as f:
                    contents = f.read()
                    exception_str = f"{exception_str}\n\n{contents}"
        except Exception:
            pass

    return exception_str


def format_code(source, fast, configuration):
    try:
        mode = black.FileMode(**configuration)
        formatted = black.format_file_contents(source, fast=fast, mode=mode)
    except black.NothingChanged:
        formatted = source
    except Exception as exc:
        formatted = normalize_exception(exc)
        return (formatted, 400)

    return formatted


@app.route("/", methods=["POST", "GET"])
@cross_origin()
def index():
    if request.method == "POST":
        data = request.get_json()
        source = data.get("source")
        options = data.get("config", {})
        line_length = int(options.get("line_length", 88))
        skip_string_normalization = bool(
            options.get("skip_string_normalization", False)
        )
        skip_magic_trailing_comma = bool(
            options.get("skip_magic_trailing_comma", False)
        )
        preview = bool(options.get("preview", False))
        py36 = bool(options.get("py36", False))
        pyi = bool(options.get("pyi", False))
        fast = bool(options.get("fast", False))
        target_versions = options.get("target_versions", set())

    else:
        state = request.args.get("state")

        if state:
            state = decompress_state(state)
            line_length = int(state.get("ll", 88))
            skip_string_normalization = bool(state.get("ssn", False))
            skip_magic_trailing_comma = bool(state.get("smtc", False))
            preview = bool(state.get("prv", False))
            py36 = bool(state.get("py36", False))
            pyi = bool(state.get("pyi", False))
            fast = bool(state.get("fast", False))
            target_versions = state.get("tv", set())
        else:
            line_length = 88
            skip_string_normalization = False
            skip_magic_trailing_comma = False
            py36 = False
            pyi = False
            fast = False
            preview = False
            target_versions = set()

    if py36:
        target_versions.add(TargetVersion.PY36)

    if request.method == "POST":
        formatted = format_code(
            source,
            fast=fast,
            configuration={
                "target_versions": {TARGET_VERSIONS[t] for t in target_versions},
                "line_length": line_length,
                "is_pyi": pyi,
                "string_normalization": not skip_string_normalization,
                "magic_trailing_comma": not skip_magic_trailing_comma,
                "preview": preview,
            },
        )
        return formatted

    if BLACK_VERSION == "stable":
        version = f"v{black_version}"
    else:
        version = f"https://github.com/psf/black/commit/{black_version}"

    return jsonify({
        "options": {
            "line_length": line_length,
            "skip_string_normalization": skip_string_normalization,
            "skip_magic_trailing_comma": skip_magic_trailing_comma,
            "target_versions": list(target_versions),
            "pyi": pyi,
            "fast": fast,
            "preview": preview,
        },
        "version": black_version,
    })


@app.route("/version", methods=["GET"])
@cross_origin()
def version():
    return jsonify({"version": black_version})


if __name__ == "__main__":
    app.run(debug=True)
