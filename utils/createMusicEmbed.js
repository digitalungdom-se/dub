module.exports = function createMusicEmbed( metadata ) {
  return {
    'content': '',
    'title': metadata.title,
    'description': '',
    'url': metadata.video_url,
    'color': 4086462,
    'timestamp': Date(),
    'thumbnail': {
      'url': metadata.thumbnail_url
    },
    'author': {
      'name': metadata.author.name,
      'url': metadata.author.user_url,
      'icon_url': metadata.author.avatar
    }
  };
}