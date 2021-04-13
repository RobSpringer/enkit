def labelrelative(label, path):
    """Computes a path relative to the root of a label.

    Bazel saves files in a different directory for each target.

    For example, if you have a target named "dplib" under "enkit/parser",
    an output file called "lib/amd/dplib.so" will likely be saved as
    "enkit/parser/dpilib/lib/amd/dplib.so".

    If you reference this file from another rule, you will see this
    path in the short_path attribute of the corresponding File object.

    This function allows to compute the original path of the output
    file: it removes the prefix the original target added. In the
    example above, given a Label() object representing the ":dpilib"
    target and the short_path of a file, it will return the relative
    path to that label, "lib/amd/dpilib.so".

    Args:
      label: a Label() object from bazel.
      path: a string, representing the name of a file generated by that label.
    Returns:
      string, path relative to where all the output files for that
      label are saved by bazel.
    """
    to_strip = "%s/%s" % (label.package, label.name)
    if path == to_strip:
        return ""

    if path.startswith(to_strip + "/"):
        path = path[len(to_strip) + 1:]
    return path