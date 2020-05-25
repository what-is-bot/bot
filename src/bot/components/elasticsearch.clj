(ns bot.components.elasticsearch
  (:require [bot.protocols.config :as c-pro]
            [bot.protocols.db :as db-pro]
            [bot.utils :as utils]
            [clojure.tools.logging :refer [debug error]]
            [com.stuartsierra.component :as component]
            [qbits.spandex :as sp]))

(defonce ^:private index :entries)
(defonce ^:private index-mapping {:properties {:question {:type :text}
                                               :answers  {:type       :nested
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

(defn- insert*! [client document]
  (debug "Inserting..." (sp/request client {:method :post
                                            :url    [index :_doc]
                                            :body   document})))

(defn- update*! [client {:keys [id _id] :as document}]
  {:pre [(some? (or id, _id))]}
  (let [id (or id _id)]
    (debug "Updating..." (sp/request client {:method :post
                                             :url    [index :_update id]
                                             :body   (dissoc document :id :_id)}))))

(defrecord ElasticSearch [config client]
  db-pro/Db
  (insert! [this document]
    (insert*! (:client this) document))
  (update! [this document]
    (update*! (:client this) document))
  (upvote! [this document-id])
  (upvote! [this document-id qty])
  (downvote! [this document-id])
  (downvote! [this document-id qty])
  (query [this question])

  component/Lifecycle
  (start [this]
    (let [hosts  (c-pro/get! config [:elasticsearch :hosts])
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
