export const VERSION_MAJOR = 1;
export const VERSION_MINOR = 2;
export const VERSION_PATCH = 9;
export const VERSION_PRE = "202605021353";

export const formatVersion = (major, minor, patch, pre) => {
  return `${major}.${minor}.${patch}-${pre}`;
};

export const versionInfo = {
  version: formatVersion(
    VERSION_MAJOR,
    VERSION_MINOR,
    VERSION_PATCH,
    VERSION_PRE,
  ),
};
