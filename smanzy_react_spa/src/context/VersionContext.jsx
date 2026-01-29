import { createContext, useContext } from "react";

export const VersionContext = createContext({
  versionInfo: null,
  isLoading: true,
  error: null,
});

export const useVersion = () => {
  const context = useContext(VersionContext);
  if (!context) {
    throw new Error("useVersion must be used within a VersionProvider");
  }
  return context;
};
