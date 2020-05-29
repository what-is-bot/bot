(ns bot.components.elasticsearch
  (:require [bot.protocols.config :as protocols.config]
            [bot.protocols.faq :as protocols.faq]
            [bot.schema.question :as schema.question]
            [clojure.tools.logging :refer [debug error]]
            [com.stuartsierra.component :as component]
            [qbits.spandex :as sp]
            [schema.core :as s])
  (:import org.elasticsearch.client.RestClient))

(defonce ^:private index :questions)
;; research about parent/child mapping
;; https://www.elastic.co/blog/managing-relations-inside-elasticsearch
(defonce ^:private index-mapping {:properties {:text        {:type :text}
                                               :asked-by-id {:type  :keyword
                                                             :index false}
                                               :answers     {:type       :nested
                                                             :properties {:text      {:type :text}
                                                                          :upvotes   {:type :integer}
                                                                          :downvotes {:type :integer}
                                                                          :score     {:type :integer}}}}})
(defn- index-exist? [client]
  (try
    (sp/request client {:method :get
                        :url [index]})
    true
    (catch Exception ex
      (error (.getMessage ex) (ex-data ex))
      false)))

(defn- setup-indexes! [client]
  (when-not (index-exist? client)
    (debug "Creating indexes...")
    (debug "Response" (sp/request client {:method :put
                                          :url    [index]
                                          :body   {:mappings index-mapping}}))))

(s/defn ^:private ask :- schema.question/Question
  [client :- RestClient
   question :- schema.question/Question])

(s/defrecord ElasticSearch [config :- protocols.config/Config
                            client :- RestClient]
  protocols.faq/Faq
  (ask! [this question])
  (query [this & {:keys [text asked-by answered-by] :as args}])
  (answer! [this question answer])
  (upvote! [this answer who])
  (upvote! [this answer who delta])
  (downvote! [this answer who])
  (downvote! [this answer who delta])

  component/Lifecycle
  (start [this]
    (let [hosts  (protocols.config/get! config [:elasticsearch :hosts])
          client (sp/client {:hosts hosts})
          component (assoc this :client client)]
      (setup-indexes! client)
      component))
  (stop [this])

  Object
  (toString [_] "<ElasticSearch"))

(defmethod print-method ElasticSearch [_ ^java.io.Writer w]
  (.write w "<ElasticSearch>"))

(defn new-elastic-search [] (map->ElasticSearch {}))
