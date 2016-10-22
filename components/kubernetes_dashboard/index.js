$(function() {
    window.custom = {
        serverFromQueryString: '',
        namespaceFromQueryString: '',

        getParameterByName: function(name, url) {
            if (!url) url = window.location.href;
            name = name.replace(/[\[\]]/g, "\\$&");
            var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
                results = regex.exec(url);

            if (!results) return ''; // the name not exists
            if (!results[2]) return '';

            return decodeURIComponent(results[2].replace(/\+/g, " "));
        },

        objectToQueryString: function(obj) {
            var s = '';
            for (var key in obj) {
                if (obj.hasOwnProperty(key)) {
                    s += '&' + encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]);
                }
            }

            if (s.length && s[0] === '&') {
                s = s.substring(1);
            }

            return s;
        },

        // Describe job and show pods of it
        describeJobAndPods: function(btn) {
            var $row = $(btn).closest('tr'),
                param = {
                    action: 'describe_job',
                    name: $row.data('name'),
                    namespace: $row.data('namespace'),
                    kubernetes_api_server: this.serverFromQueryString
                };

            window.open('/kubernetes_dashboard/?'+this.objectToQueryString(param));
        },

        // Describe pod and show logs
        describePodAndLogs: function(btn) {
            var $row = $(btn).closest('tr'),
                param = {
                    action: 'describe_pod',
                    name: $row.data('name'),
                    namespace: $row.data('namespace'),
                    kubernetes_api_server: this.serverFromQueryString
                };

            window.open('/kubernetes_dashboard/?'+this.objectToQueryString(param));
        },

        // Delete a job
        deleteJob: function(btn, noConfirm, fails) {
            var $row = $(btn).closest('tr');

            if (!noConfirm) {
                if (!confirm("Sure?")) { return; }
            }

            $.post('/kubernetes_dashboard/delete_job', {
                name: $row.data('name'),
                namespace: $row.data('namespace'),
                kubernetes_api_server: this.serverFromQueryString
            }, function(data) {
                if (data.Error) {
                    if (typeof fails === 'function') {
                        fails(data.Error);
                    } else {
                        alert(data.Error);
                    }

                    return;
                }

                $row.fadeOut(300, function() { $row.remove(); });
            });
        },

        // Delete a pod
        deletePod: function(btn) {
            var $row = $(btn).closest('tr');

            if (!confirm("Sure?")) { return; }

            $.post('/kubernetes_dashboard/delete_pod', {
                name: $row.data('name'),
                namespace: $row.data('namespace'),
                kubernetes_api_server: this.serverFromQueryString
            }, function(data) {
                if (data.Error) { alert(data.Error); return; }
                $row.remove();
            });
        },

        // show panel of Describe Job
        showPanelOfDescribeJob: function() {
            var $panel = $('#job-panel'),
                $body = $panel.find('tbody'),
                $row = $body.find('tr').detach(),
                param = {
                    name: this.getParameterByName('name'),
                    namespace: this.namespaceFromQueryString,
                    kubernetes_api_server: this.serverFromQueryString
                };

            $panel.show();

            $panel.find('.delete-job-btn').click(function() {
                if (!confirm("Sure?")) { return; }

                $.post('/kubernetes_dashboard/delete_job', param, function(data) {
                    if (data.Error) { alert(data.Error); return; }
                    $panel.html('<div class="empty">Job deleted.</div>');
                });
            });

            $.post('/kubernetes_dashboard/describe_job', param, function(data) {
                if (data.Error) { alert(data.Error); return; }
                if (!data.Data) { alert('No data.'); return; }

                $panel.find('.toolbar').show();
                $panel.find('.head-tip code').text(param.name);
                $panel.find('pre code').text(data.Data.Job);

                // process pods
                if (!data.Data.PodList || !data.Data.PodList.items || data.Data.PodList.items.length===0) {
                    $body.append('<tr><td colspan="5" class="empty">No pod.</td></tr>');
                    return;
                }

                data.Data.PodList.items.sort(function(i, j) {
                    return i.metadata.creationTimestamp < j.metadata.creationTimestamp ? 1 : -1;
                });
                data.Data.PodList.items.forEach(function(pod, i) {
                    var $r = $row.clone();
                    $r.attr('data-name', pod.metadata.name);
                    $r.attr('data-namespace', pod.metadata.namespace);
                    $r.find('.order').text('#'+(i+1));
                    $r.find('.name').text(pod.metadata.name);
                    $r.find('.creationTimestamp').text(pod.metadata.creationTimestamp);
                    $r.find('.node').text(pod.spec.nodeName);
                    $r.appendTo($body);
                });
            });
        },

        showPanelOfDescribePod: function() {
            var $panel = $('#pod-panel'),
                param = {
                    name: this.getParameterByName('name'),
                    namespace: this.namespaceFromQueryString,
                    kubernetes_api_server: this.serverFromQueryString
                };

            $panel.show();

            $panel.find('.delete-pod-btn').click(function() {
                if (!confirm("Sure?")) { return; }

                $.post('/kubernetes_dashboard/delete_pod', param, function(data) {
                    if (data.Error) { alert(data.Error); return; }
                    $panel.html('<div class="empty">Pod deleted.</div>');
                });
            });

            $.post('/kubernetes_dashboard/describe_pod', param, function(data) {
                if (data.Error) { alert(data.Error); return; }
                if (!data.Data) { alert('No data.'); return; }

                $panel.find('.toolbar').show();
                $panel.find('.head-tip code').text(param.name);
                $('#pod-description').text(data.Data.Pod);
                $('#pod-logs pre code').text(data.Data.Logs);
            });
        },

        showPanelOfJobsSearchResult: function() {
            var $panel = $('#jobs-panel'),
                $toolbar = $panel.find('.toolbar'),
                $cancel = $toolbar.find('.cancel'),
                $result = $toolbar.find('.result'),
                $body = $panel.find('tbody'),
                $row = $body.find('tr').detach(),
                interval = null,
                wantStopDeleting = false,
                param = {
                    namespace: this.namespaceFromQueryString,
                    kubernetes_api_server: this.serverFromQueryString
                };

            $panel.show();

            $cancel.click(function() {
                wantStopDeleting = true;
                $cancel.hide();
                $result.text('Deleting cancel.').show();
            });

            $toolbar.find('.delete-jobs-btn').click(function() {
                if (!confirm("Sure?")) { return; }

                var $rows = $body.find('tr'),
                    $i = $rows.length - 1;

                $(this).hide();
                $cancel.show();

                interval = setInterval(function() {
                    if ($i < 0 || wantStopDeleting) {
                        clearInterval(interval);
                        interval = null;
                        
                        if ($i < 0) {
                            $cancel.hide();
                            $result.text('Deleting finished.').show();
                        }

                        return;
                    }

                    var $r = $rows.eq($i);

                    $i -= 1;

                    $cancel.text($r.find('.order').text()+' job is deleting... Click me to cancel');
                    custom.deleteJob($r.find('.delete'), true, function(err) {
                        // append error info if fails
                        $r.find('.delete').closest('td').append(err);
                    });
                }, 1000); // wait 1 second for possible to cancel
            });

            $.post('/kubernetes_dashboard/get_jobs', param, function(data) {
                if (data.Error) { alert(data.Error); return; }
                if (!data.Data || !Array.isArray(data.Data.items)) { alert('No data.'); return; }

                var jobs = data.Data.items;

                $panel.find('.head-tip code').text(jobs.length);
                if (jobs.length===0) {
                    $body.append('<tr><td colspan="6" class="empty">No job.</td></tr>');
                    return;
                }

                jobs.sort(function(i, j) {
                    return i.metadata.creationTimestamp < j.metadata.creationTimestamp ? 1 : -1;
                });
                jobs.forEach(function(job, i) {
                    var $r = $row.clone();
                    $r.attr('data-name', job.metadata.name);
                    $r.attr('data-namespace', job.metadata.namespace);
                    $r.find('.order').text('#'+(i+1));
                    $r.find('.namespace').text(job.metadata.namespace);
                    $r.find('.name').text(job.metadata.name);
                    $r.find('.creationTimestamp').text(job.metadata.creationTimestamp);
                    $r.find('.succeeded').text(job.status.succeeded || 0);
                    $r.appendTo($body);
                });

                $toolbar.show();
            });
        },

        loadPage: function() {
            var $server = $('input[name="kubernetes_api_server"]'),
                $namespace = $('input[name="namespace"]');

            this.serverFromQueryString = this.getParameterByName('kubernetes_api_server');
            this.namespaceFromQueryString = this.getParameterByName('namespace');
            $server.val(this.serverFromQueryString);
            $namespace.val(this.namespaceFromQueryString);

            $server.focus();

            $('#search-jobs').click(function() {
                var server = $server.val().trim();
                if (server.length === 0) { alert('kubernetes API Server must be not empty'); return; }
                if (server.toLowerCase().indexOf('://') === -1) { $server.val('http://'+server); }
                location.href = '/kubernetes_dashboard/?action=job_search_result&'+$('form').serialize();
            });

            switch (this.getParameterByName('action')) {
                case '':
                    break;
                case 'describe_job':
                    this.showPanelOfDescribeJob();
                    break;
                case 'describe_pod':
                    this.showPanelOfDescribePod();
                    break;
                case 'job_search_result':
                    this.showPanelOfJobsSearchResult();
                    break;
                default:
                    alert('Unkown action');
            }
        }
    };

    custom.loadPage();
});
