/* global db */

module.exports.getUserByDiscordId = async function createUser( discordId ) {
  return db.collection( 'users' ).findOne( { 'connectedApps.discord': discordId } );
};