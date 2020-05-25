(ns bot.components.config
  (:require [bot.protocols.config :as c-pro]
            [com.stuartsierra.component :as component]
            [config.core :refer [load-env]]))

(defn- do-get [this path]
  (-> this :config (get-in path)))

(defrecord Config [config]
  c-pro/Config
  (get! [this path]
    (or (do-get this path)
        (throw (ex-info "Config not found" {:path path}))))
  (get* [this path]
    (do-get this path))

  component/Lifecycle
  (start [this]
    (assoc this :config (load-env)))
  (stop [this]
    (dissoc this :config))

  Object
  (toString [_] "<Config>"))

(defmethod print-method Config [_ ^java.io.Writer w]
  (.write w "<Config>"))

(defn new-config [] (map->Config {}))
