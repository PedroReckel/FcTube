import { LikeButton } from "./LikeButton";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

// Antes da versão 15 do Next:
  // Fetch - cache por padrão

// Depois da versão 15 do Next:
  // 15  Segundos por padrão

export async function getLikes(videoId: number): Promise<number> {
  await sleep(2000);
  const response = await fetch(`${process.env.DJANGO_API_URL}/videos/${videoId}/likes`, {
    next: {
      revalidate: 60,
    },
  });

  return (await response.json()).likes;
}

export type VideoLikeCounterProps = {
  videoId: number;
  likes?: number;
};

export async function VideoLikeCounter(props: VideoLikeCounterProps) {
  const { videoId, likes: propLikes } = props;
  const likes = propLikes ? propLikes : await getLikes(videoId);
  return <LikeButton videoId={videoId} likes={likes} />;
}