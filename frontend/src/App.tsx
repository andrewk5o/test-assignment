import { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { useCityCount } from "@/hooks/useCityCount";

const LETTER_RE = /^\p{L}?$/u;

export default function App() {
  const [letter, setLetter] = useState("");
  const { data, error, isFetching } = useCityCount(letter);

  function onChange(e: React.ChangeEvent<HTMLInputElement>) {
    const value = e.target.value.trim();
    const last = value.length > 1 ? [...value].at(-1)! : value;
    if (LETTER_RE.test(last)) setLetter(last);
  }

  return (
    <main className="min-h-svh bg-background text-foreground flex items-start justify-center px-4 py-16">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Cities by letter</CardTitle>
          <CardDescription>
            Top 50 cities by population (GeoNames)
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-center gap-4">
            <Input
              aria-label="Letter"
              value={letter}
              onChange={onChange}
              placeholder="?"
              autoFocus
              autoComplete="off"
              className="size-14 shrink-0 text-center font-mono text-2xl uppercase placeholder:normal-case"
            />
            <div aria-live="polite" className="min-w-0 flex-1">
              {letter === "" ? (
                <p className="text-sm text-muted-foreground">Type a letter</p>
              ) : isFetching && !data ? (
                <Skeleton className="h-8 w-20" />
              ) : error ? (
                <p role="alert" className="text-sm text-destructive">
                  {error.message}
                </p>
              ) : data ? (
                <p className="text-3xl font-semibold tabular-nums">
                  {data.count}
                  <span className="ml-1.5 text-sm font-normal text-muted-foreground">
                    {data.count === 1 ? "city" : "cities"}
                  </span>
                </p>
              ) : null}
            </div>
          </div>
          {data && letter !== "" && data.cities.length > 0 && (
            <p className="text-sm text-muted-foreground">
              {data.cities.join(", ")}
            </p>
          )}
        </CardContent>
      </Card>
    </main>
  );
}
