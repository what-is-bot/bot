(defproject bot "0.1.0-SNAPSHOT"
  :description "What is bot, clojure impl"
  :url "http://example.com/FIXME"
  :license {:name "EPL-2.0 OR GPL-2.0-or-later WITH Classpath-exception-2.0"
            :url  "https://www.eclipse.org/legal/epl-2.0/"}

  :dependencies [[org.clojure/clojure "1.10.0"]
                 [org.clojure/tools.logging "1.1.0"]
                 [ch.qos.logback/logback-classic "1.1.3"]
                 [prismatic/schema "1.1.12"]

                 [com.stuartsierra/component "1.0.0"]
                 [yogthos/config "1.1.7"]
                 [cc.qbits/spandex "0.7.4"]]

  :target-path "target/%s"

  :profiles {:uberjar {:aot :all}
             :prod    {:resource-paths ["config/prod"]}
             :test    {:dependencies [[clojure.test "1.10.1"]]}
             :dev     {:resource-paths ["config/dev"]
                       :injections     [(require '[bot.components :as c])]
                       :dependencies   [[clj-kondo "2020.04.05"]]
                       :repl-options   {:init-ns bot.components}}

             :repl-start {:repl-options {:prompt  #(str "[bot] " % "=> ")
                                         :timeout 300000
                                         :init-ns bot.components}}}

  :aliases {"clj-kondo" ["run" "-m" "clj-kondo.main" "--config" ".clj-kondo/config.edn" "--lint" "src" "test"]
            "lint"      ["do" ["clj-kondo"]]}

  :main ^{:skip-aot false} bot.core)
