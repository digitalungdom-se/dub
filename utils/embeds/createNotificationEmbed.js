module.exports = function createNotificationEmbed( type, message, colour, author ) {
  return {
    'title': type,
    'description': `<@${author.id}>, ${message}`,
    'color': colour,
    'timestamp': Date(),
    'author': {
      'name': author.name,
      'icon_url': author.url
    }
  };
};