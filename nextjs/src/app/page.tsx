
import VideoCardSkeleton from "@/components/VideoCardSkeleton";
import { VideosList } from "@/components/VideosList";
import { Suspense } from "react";

// const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

// async function getVideos(): Promise<VideoModel[]> {
//   const response = await fetch('http://localhost:8000/api/videos', {
//     // next: {
//     //   revalidate: 10, // Eu posso ter um chache durante 10 segungos (diminuir as chamadas na API)
//     // },
//     cache: "no-store" // Não guardar cache (buscar tudo no servidor)
//   });
//   await sleep(2000) // Tornar a chamada na API mais lenta
//   return response.json();
// }

export default async function Home({searchParams}: {searchParams: {search: string}}) {
  return (
    <div className="container mx-auto px-4 py-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <Suspense
          fallback={new Array(15).fill(null).map((_, index) => (
            <VideoCardSkeleton key={index} />
          ))}

        >
          <VideosList search={searchParams.search}/>
        </Suspense>
      </div>
    </div>
  );
}
