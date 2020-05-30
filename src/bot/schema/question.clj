(ns bot.schema.question
  (:require [bot.schema.answer :as schema.answer]
            [bot.schema.user :as schema.user]
            [schema.core :as s]))

(s/defschema Question {(s/optional-key :id)         s/Uuid
                       (s/optional-key :created-at) s/Inst
                       (s/optional-key :asked-by)   schema.user/User
                       :text                        s/Str
                       (s/optional-key :answers)    [schema.answer/Answer]})
