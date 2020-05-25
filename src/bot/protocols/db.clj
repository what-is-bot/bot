(ns bot.protocols.db)

(defprotocol Db
  (insert! [this document])
  (update! [this document])
  (upvote! [this document-id] [this document-id qty])
  (downvote! [this document-id] [this document-id qty])
  (query [this question]))
