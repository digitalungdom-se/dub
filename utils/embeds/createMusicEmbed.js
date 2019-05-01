module.exports = function createMusicEmbed( metadata ) {
  let seconds = metadata.length_seconds % 60;
  if ( seconds < 10 ) seconds = `0${seconds}`;
  const minutes = Math.floor( metadata.length_seconds / 60 );

  return {
    'title': metadata.title,
    'description': `Spelas nu på begäran av <@${metadata.user}>. Den är ${minutes}:${seconds} minuter lång.`,
    'url': metadata.video_url,
    'color': 4086462,
    'timestamp': ( new Date() ).toISOString(),
    'thumbnail': {
      'url': metadata.thumbnail_url
    },
    'author': {
      'name': metadata.author.name,
      'url': metadata.author.user_url,
      'icon_url': metadata.author.avatar
    }
  };
};