// get url by env, if dev, return xxx:8081, if prod, return xxx:8080
const getEnvUrl = () => {
  if (process.env.NODE_ENV === "development") {
    return "http://localhost:8080";
  } else {
    return window.location.origin;
  }
};

const getCurrDirPath = (): string => {
  const url = new URL(window.location.href);
  const dirPath = url.searchParams.get("path") || "";

  return dirPath;
};

const joinPath = (...parts: string[]): string => {
  return parts.reduce((previous, current) => {
    if (previous.endsWith("/")) {
      previous = previous.slice(0, -1);
    }
    if (current.startsWith("/")) {
      current = current.slice(1);
    }
    return `${previous}/${current}`;
  });
};

export { getEnvUrl, getCurrDirPath, joinPath };
