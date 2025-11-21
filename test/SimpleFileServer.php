<?php

class SimpleFileServer
{
    private $baseUrl;
    private $apiToken;

    public function __construct($baseUrl, $apiToken)
    {
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->apiToken = $apiToken;
    }

    /**
     * Ping the server to check status.
     *
     * @return array Response from the server.
     */
    public function ping()
    {
        return $this->sendGetRequest('/_admin/ping');
    }

    /**
     * Upload a single file.
     *
     * @param string $filePath Remote file path to save.
     * @param string $localFilePath Local file path to upload.
     * @return array Response from the server.
     */
    public function upload($filePath, $localFilePath)
    {
        if (!file_exists($localFilePath)) {
            throw new Exception("Local file does not exist: $localFilePath");
        }

        $postData = [
            'file' => new CURLFile($localFilePath),
            'filePath' => $filePath,
        ];

        return $this->sendPostRequest('/_admin/upload', $postData);
    }

    /**
     * Upload content
     */
    public function uploadContent($filePath, $content)
    {
        if (empty($filePath)) {
            throw new Exception("File path cannot be empty");
        }
        $tempFile = tempnam(sys_get_temp_dir(), 'upload_');
        file_put_contents($tempFile, $content);
        $ret = $this->upload($filePath, $tempFile);
        unlink($tempFile); // Clean up temporary file
        return $ret;
    }

    /**
     * Initialize multipart upload.
     *
     * @param string $filePath Remote file path.
     * @param int $totalParts Total number of parts.
     * @param int $totalSize Total file size in bytes.
     * @return array Response from the server, including uploadId.
     */
    public function initMultipartUpload($filePath, $totalParts, $totalSize)
    {
        $data = [
            'filePath' => $filePath,
            'totalParts' => $totalParts,
            'totalSize' => $totalSize,
        ];

        return $this->sendJsonPostRequest('/_admin/upload/multipart_init', $data);
    }

    /**
     * Upload a part of the file.
     *
     * @param string $uploadId Upload ID from init.
     * @param int $partNumber Part number.
     * @param string $partFilePath Local path to the part file.
     * @return array Response from the server.
     */
    public function uploadPart($uploadId, $partNumber, $partFilePath)
    {
        if (!file_exists($partFilePath)) {
            throw new Exception("Part file does not exist: $partFilePath");
        }

        $postData = [
            'uploadId' => $uploadId,
            'partNumber' => $partNumber,
            'file' => new CURLFile($partFilePath),
        ];

        return $this->sendPostRequest('/_admin/upload/multipart_upload', $postData);
    }

    /**
     * Complete multipart upload.
     *
     * @param string $uploadId Upload ID.
     * @return array Response from the server.
     */
    public function completeMultipartUpload($uploadId)
    {
        $data = [
            'uploadId' => $uploadId,
        ];

        return $this->sendJsonPostRequest('/_admin/upload/multipart_end', $data);
    }

    /**
     * Abort multipart upload.
     *
     * @param string $uploadId Upload ID.
     * @return array Response from the server.
     */
    public function abortMultipartUpload($uploadId)
    {
        $data = [
            'uploadId' => $uploadId,
        ];

        return $this->sendJsonPostRequest('/_admin/upload/abort', $data);
    }

    /**
     * Check if a file exists.
     *
     * @param string $fileName File name to check.
     * @return bool True if file exists.
     */
    public function has($fileName)
    {
        $data = ['path' => $fileName];
        $response = $this->sendJsonPostRequest('/_admin/has', $data);
        return isset($response['data']) && $response['data'] === true;
    }

    /**
     * Get the size of a file.
     *
     * @param string $fileName File name.
     * @return int File size in bytes, or -1 on error.
     */
    public function size($fileName)
    {
        $data = ['path' => $fileName];
        $response = $this->sendJsonPostRequest('/_admin/size', $data);
        return isset($response['data']['size']) ? $response['data']['size'] : -1;
    }

    /**
     * Get the content of a file.
     *
     * @param string $fileName File name.
     * @return string|false File content or false on error.
     */
    public function get($fileName)
    {
        $url = $this->baseUrl . '/_admin/get';
        $ch = curl_init($url);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode(['path' => $fileName]));
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'admin-api-token: ' . $this->apiToken,
            'Content-Type: application/json',
        ]);
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        if ($httpCode == 200) {
            return $response;
        } else {
            return false;
        }
    }

    /**
     * Move a file.
     *
     * @param string $from Source path.
     * @param string $to Destination path.
     * @return array Response from the server.
     */
    public function move($from, $to)
    {
        $data = [
            'from' => $from,
            'to' => $to,
        ];
        return $this->sendJsonPostRequest('/_admin/move', $data);
    }

    /**
     * Delete a file.
     *
     * @param string $path File path to delete.
     * @return array Response from the server.
     */
    public function delete($path)
    {
        $data = [
            'path' => $path,
        ];
        return $this->sendJsonPostRequest('/_admin/delete', $data);
    }

    private function sendGetRequest($endpoint)
    {
        $url = $this->baseUrl . $endpoint;
        $ch = curl_init($url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        return $this->parseResponse($response, $httpCode);
    }

    private function sendPostRequest($endpoint, $postData)
    {
        $url = $this->baseUrl . $endpoint;
        $ch = curl_init($url);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $postData);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'admin-api-token: ' . $this->apiToken,
        ]);
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        return $this->parseResponse($response, $httpCode);
    }

    private function sendJsonPostRequest($endpoint, $data)
    {
        $url = $this->baseUrl . $endpoint;
        $ch = curl_init($url);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'admin-api-token: ' . $this->apiToken,
            'Content-Type: application/json',
        ]);
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        return $this->parseResponse($response, $httpCode);
    }

    private function parseResponse($response, $httpCode)
    {
        if ($httpCode != 200) {
            return [
                'code' => -1,
                'msg' => 'http code ' . $httpCode,
                'data' => [
                    '_raw' => $response,
                ]
            ];
        }

        $data = json_decode($response, true);
        if (json_last_error() !== JSON_ERROR_NONE) {
            return [
                'code' => -1,
                'msg' => 'Invalid JSON response',
                'data' => [
                    '_raw' => $response,
                ]
            ];
        }
        return $data;
    }
}


// $client = new SimpleFileServer(
//     'http://127.0.0.1:60088',
//     'xxx'
// );

// upload content
// $ret = $client->uploadContent('data/mp4.mp4', file_get_contents('/Users/mz/data/material/mp4.mp4'));
// echo json_encode($ret) . PHP_EOL;

// // Test Ping
// $pingResult = $client->ping();
// echo json_encode($pingResult) . PHP_EOL;

// // Test Single File Upload
// echo "Testing Single File Upload:\n";
// $testFile = '/Users/mz/data/material/mp4.mp4'; // Assume this file exists
// $uploadResult = $client->upload('data/mp4.mp4', $testFile);
// echo json_encode($uploadResult) . PHP_EOL;

// // Test Multipart Upload
// echo "Testing Multipart Upload:\n";
// $largeFile = '/Users/mz/data/material/mp4.mp4';

// $content = file_get_contents($largeFile);
// $partSize = 100000; // 100KB per part
// $totalSize = filesize($largeFile);
// $totalParts = ceil($totalSize / $partSize);

// $initResult = $client->initMultipartUpload('data/mp4-large.mp4', $totalParts, $totalSize);
// echo json_encode($initResult) . PHP_EOL;
// $uploadId = $initResult['data']['uploadId'];
// for ($i = 1; $i <= $totalParts; $i++) {
//     $partContent = substr($content, ($i - 1) * $partSize, $partSize);
//     $partFile = "/tmp/part_$i.mp4";
//     file_put_contents($partFile, $partContent);
//     $partResult = $client->uploadPart($uploadId, $i, $partFile);
//     echo json_encode($partResult) . PHP_EOL;
//     unlink($partFile); // Clean up
//     echo "upload $i\n";
// }
// $completeResult = $client->completeMultipartUpload($uploadId);
// echo json_encode($completeResult) . PHP_EOL;
// echo "\n";

// // Test Has File
// echo "Testing Has File:\n";
// $hasResult = $client->has('data/mp4.mp4');
// echo "Has file: " . ($hasResult ? 'Yes' : 'No') . PHP_EOL;

// // Test Size File
// echo "Testing Size File:\n";
// $sizeResult = $client->size('data/mp4.mp4');
// echo "File size: " . $sizeResult . " bytes" . PHP_EOL;

// // Test Get File Content (if small file)
// echo "Testing Get File Content:\n";
// $getResult = $client->get('data/mp4.mp4');
// if ($getResult !== false) {
//     echo "File content length: " . strlen($getResult) . " bytes" . PHP_EOL;
// } else {
//     echo "Failed to get file content" . PHP_EOL;
// }

// // Test Move File
// echo "Testing Move File:\n";
// $moveResult = $client->move('data/mp4.mp4', 'data/moved_mp4.mp4');
// echo json_encode($moveResult) . PHP_EOL;

// // Test Delete File
// echo "Testing Delete File:\n";
// $deleteResult = $client->delete('data/moved_mp4.mp4');
// echo json_encode($deleteResult) . PHP_EOL;

// Test Multipart Download or other methods if needed
// Note: Assuming the server supports these operations
