(ns bot.components
  (:require [bot.components.config :as config]
            [bot.system :as system]
            [com.stuartsierra.component :as component]
            [bot.components.elasticsearch :as elasticsearch]
            [clojure.tools.logging :refer [debug]]
            [bot.protocols.faq :as p-faq]))

(defn system-map []
  (component/system-map
   :config (config/new-config)
   :db     (component/using (elasticsearch/new-elastic-search) [:config])))

(defn create-and-start-system! []
  (debug "Starting system...")
  (system/bootstrap! (system-map)))

(defn restart-system! []
  (system/destroy!)
  (create-and-start-system!))

(defn test! []
  (let [sys               (restart-system!)
        es                (:db sys)
        question          (p-faq/ask! es {:text "questão"})
        answer            {:text "resposta"}]
    (p-faq/answer! es question answer)))
