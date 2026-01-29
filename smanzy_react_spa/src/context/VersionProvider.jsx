import { useQuery } from "@tanstack/react-query";
import api from "@/services/api";
import { VersionContext } from "./VersionContext";

export const VersionProvider = ({ children }) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["version"],
    queryFn: () => api.get("/version").then((res) => res.data),
    staleTime: Infinity, // Version info shouldn't change during a session
    retry: 2,
  });

  return (
    <VersionContext.Provider value={{ versionInfo: data, isLoading, error }}>
      {children}
    </VersionContext.Provider>
  );
};
