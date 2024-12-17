import { VideoModel } from "../../../models";

export async function getVideo(slug: string): Promise<VideoModel> {
  const response = await fetch(`${process.env.DJANGO_API_URL}/videos/${slug}`, {
    // cache: "no-cache",
    next: {
        // revalidate: 60 * 5 // Aqui ele só vai fazer cache por 5 minutos
        tags: [`video-${slug}`] // Aqui o cache seria eterno até o vídeo ser modificado
    }
  });
  return response.json();
}

//revalidate on demand
// /admin/videos -> django -> http -> next.js -> revalidate
