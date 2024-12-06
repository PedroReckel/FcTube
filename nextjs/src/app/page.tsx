import { VideoCard } from "@/components/VideoCard";

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
        <VideoCard title="Video Title" thumbnail="/thumbnail.jpg" views={1000} likes={100}/>
      </div>
    </div>
  );
}
