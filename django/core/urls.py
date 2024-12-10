from django.urls import path
from core.api import videos_detail_by_id, videos_detail_by_slug, videos_list, videos_list_recommended

urlpatterns = [
    path('api/videos', videos_list, name='api_videos_list'),
    path('api/videos/<int:id>', videos_detail_by_id, name='videos_detail_by_id'),
    path('api/videos/<slug:slug>', videos_detail_by_slug, name='video_detail_by_slug'),
    path('api/videos/<int:id>/recommended', videos_list_recommended, name='api_videos_list_recommended'),
]