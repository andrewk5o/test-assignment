export interface CityCountResponse {
  letter: string;
  count: number;
  cities: string[];
}

export async function fetchCityCount(
  letter: string,
  signal?: AbortSignal,
): Promise<CityCountResponse> {
  const res = await fetch(
    `/api/cities/count?letter=${encodeURIComponent(letter)}`,
    { signal },
  );
  if (!res.ok) {
    let message = `Request failed (${res.status})`;
    try {
      const body = (await res.json()) as { error?: string };
      if (body.error) message = body.error;
    } catch {}
    throw new Error(message);
  }
  return res.json() as Promise<CityCountResponse>;
}
