(ns bot.protocols.config)

(defprotocol Config
  (get! [component path])
  (get* [component path]))
