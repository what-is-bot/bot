(ns bot.components.elasticsearch
  (:require [bot.protocols.config :as protocols.config]
            [bot.protocols.faq :as protocols.faq]
            [bot.schema.answer :as schema.answer]
            [bot.schema.question :as schema.question]
            [bot.time :as time]
            [clj-uuid :as uuid]
            [clojure.tools.logging :refer [error]]
            [com.stuartsierra.component :as component]
            [qbits.spandex :as sp]
            [schema.core :as s])
  (:import org.elasticsearch.client.RestClient))


;; https://www.elastic.co/guide/en/elasticsearch/reference/current/parent-join.html
(defonce join-field :answers)
(defonce parent-value :question)
(defonce child-value :answer)

(def entries [:entries {:properties {:text        {:type :text}
                                     :asked-by    {:type :keyword}
                                     :answered-by {:type :keyword}
                                     :upvotes     {:type :integer}
                                     :downvotes   {:type :integer}
                                     :score       {:type :integer}
                                     join-field   {:type      :join
                                                   :relations {parent-value child-value}}}}])

(def indexes [entries])

(defn- index-exist? [client index]
  (try
    (sp/request client {:method :get
                        :url    [index]})
    true
    (catch Exception ex
      (error (.getMessage ex) (ex-data ex))
      false)))

(defn- setup-indexes! [client]
  (doseq [[index mapping] indexes]
    (let [exists? (index-exist? client index)
          url     (if exists? [index :_mapping] [index])
          body    (if exists? mapping {:mappings mapping})]
      (sp/request client {:method       :put
                          :url  url
                          :body body}))))

(s/defn ^:private ask!* :- schema.question/Question
  [client :- RestClient
   question :- schema.question/Question]
  (let [{:keys [id] :as question} (assoc question :id (uuid/squuid) :created-at (time/inst))
        body                      (assoc question join-field parent-value)
        {:keys [status] :as resp} (sp/request client {:method :put
                                                      :url    [:entries :_doc id]
                                                      :body   body})]
    (if (and (>= status 200)
             (< status 400))
      question
      (throw (ex-info "failed to communicate with ElasticSearch" {:question question
                                                                  :resp     resp})))))

(s/defn ^:private answer!* :- schema.answer/Answer
  [client            :- RestClient
   {question-id :id} :- schema.question/Question
   answer            :- schema.answer/Answer]
  (let [{:keys [id] :as answer}   (assoc answer :id (uuid/squuid) :created-at (time/inst))
        body                      (assoc answer join-field {:name child-value :parent question-id})
        {:keys [status] :as resp} (sp/request client {:method       :put  ;; TODO maybe abstract ElasticSearch calls
                                                      :url          [:entries :_doc id]
                                                      :query-string {:routing question-id}
                                                      :body         body})]
    resp
    (if (and (>= status 200) ;; TODO better handle results
             (< status 400))
      answer
      (throw (ex-info "failed to communicate with ElasticSearch" {:answer answer
                                                                  :resp   resp})))))

(s/defrecord ElasticSearch [config :- protocols.config/Config
                            client :- RestClient]
  protocols.faq/Faq
  (ask! [this question]
    (ask!* (:client this) question))
  (query [this & {:keys [text asked-by answered-by] :as args}])
  (answer! [this question answer]
    (answer!* (:client this) question answer))
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
