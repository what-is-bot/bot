(ns bot.schema.user
  (:require [schema.core :as s]))

(s/defschema User {:id         s/Uuid
                   :created-at s/Inst
                   :username   s/Str})
