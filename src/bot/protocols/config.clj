(ns bot.protocols.config)

(defprotocol Config
  (get! [component path])
  (get* [component path]))

(def IConfig (:on-interface Config))
