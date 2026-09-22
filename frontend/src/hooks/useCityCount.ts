import { useQuery } from "@tanstack/react-query";
import { fetchCityCount } from "@/api/cities";

export function useCityCount(letter: string) {
  return useQuery({
    queryKey: ["cityCount", letter.toLowerCase()],
    queryFn: ({ signal }) => fetchCityCount(letter, signal),
    enabled: letter.length === 1,
    staleTime: 10 * 60 * 1000,
  });
}
